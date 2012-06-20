$(function() {
  var milestonesUrl = "https://api.github.com/repos/marcuswhybrow/minecraft-server-manager/milestones?callback=?",
      tagsUrl       = "https://api.github.com/repos/marcuswhybrow/minecraft-server-manager/tags?callback=?";

  $.getJSON(milestonesUrl, function(result) {
    var milestones = result.data;
    
    $.each(milestones, function(i, milestone) {
      var closedIssues = milestone.closed_issues;
      var totalIssues = closedIssues + milestone.open_issues;
      var milestonePercentage = closedIssues / totalIssues * 100;
      
      var $div = $("<div>").addClass("milestone");
      var $link =$("<a>")
            .attr("href", "https://github.com/marcuswhybrow/minecraft-server-manager/issues?state=open&milestone=" + milestone.number)
            .addClass("ga-track")
            .appendTo($div);
      $("<h2>").text(milestone.title).appendTo($link);
      $("<p>").text(milestone.description).appendTo($div);
      var $progress = $("<div>").addClass("progress").appendTo($div);
      var $bar = $("<div>").addClass("bar").css("width", milestonePercentage + "%").appendTo($progress);
      $(".milestones").append($div);
    });
    
    $(".milestones").slideDown();
  });

  function compareTags(a, b) {
    if (a[0] > b[0]) return 1;
    else if (a[0] == b[0] && a[1] > b[1]) return 1;
    else if (a[0] == b[0] && a[1] == b[1] && a[2] > b[2]) return 1;
    else if (a[0] == b[0] && a[1] == b[1] && a[2] == b[2]) return 0;
    return -1;
  }

  function matchToTagVersion(match) {
    return [ parseInt(match[1], 10), parseInt(match[2], 10), parseInt(match[3], 10) ];
  }

  $.getJSON(tagsUrl, function(result) {
    var tags = result.data,
        latestTagVersion = [0, 0, 0],
        tagVersion, match;

    $.each(tags, function(i, tag) {
      match = tag.name.match(/^(\d+)\.(\d+)\.(\d+)$/);
      if (match) {
        tagVersion = matchToTagVersion(match);
        if (compareTags(tagVersion, latestTagVersion) == 1)
          latestTagVersion = tagVersion;
      }
    });

    var $installBtnVersion = $('#install-button small');
    match = $installBtnVersion.text().match(/^\((\d+)\.(\d+)\.(\d+)\)$/);
    if (match) {
      tagVersion = matchToTagVersion(match);
      if (compareTags(latestTagVersion, tagVersion) == 1)
        $installBtnVersion.text('(' + latestTagVersion.join('.') + ')');
    } 
  });

  $('.player').tooltip();

  $('a.ga-track').live('click', function() {
    var $this = $(this);
    _gaq.push(['_trackEvent', 'External Link', $this.text(), $this.attr('href')]);
  });
});